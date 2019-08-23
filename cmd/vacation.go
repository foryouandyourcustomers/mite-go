package cmd

import (
	"errors"
	"fmt"
	"github.com/leanovate/mite-go/domain"
	"github.com/spf13/cobra"
	"strconv"
)

const (
	fullVacationDayDuration           = 8
	halfVacationDayDuration           = fullVacationDayDuration / 2
	entryListFilterAtThisYear         = "this_year"
	textProjectOrServiceNotConfigured = "please set your vacation configuration for project id, service id AND vacation days (either via arguments or config)"
	textDateFormatNotCorrect          = "format for date is not correct, expected YYYY-MM-DD, e.g. 2020-03-21"
)

var (
	vacationDetailsVerbose   bool
	vacationNote             string
	vacationHalfDay          bool
	vacationFrom             string
	vacationAmount           int
	vacationCountdownVerbose bool
)

func init() {
	vacationDetailCommand.Flags().BoolVarP(&vacationDetailsVerbose, "verbose", "v", false, "verbose output")
	vacationCommand.AddCommand(vacationDetailCommand)

	vacationCreateCommand.Flags().StringVarP(&vacationNote,  "note",    "n", "",    "A note describing your vacation")
	vacationCreateCommand.Flags().BoolVarP(&vacationHalfDay, "halfday", "d", false, "If set vacation is entered as half a day")
	vacationCreateCommand.Flags().StringVarP(&vacationFrom,  "from",    "f", "",    "Create vacation starting at date (in YYYY-MM-DD format) [Default: today]")
	vacationCreateCommand.Flags().IntVarP(&vacationAmount,   "amount",  "a", 1,     "Create amount of consecutive vacation days beginning at from date [Default: 1]")
	vacationCommand.AddCommand(vacationCreateCommand)

	vacationCountdownCommand.Flags().BoolVarP(&vacationCountdownVerbose, "verbose", "v", false, "verbose output")
    vacationCommand.AddCommand(vacationCountdownCommand)

	rootCmd.AddCommand(vacationCommand)
}

var vacationCommand = &cobra.Command{
	Use:   "vacation",
	Short: "manage your vacation",
	RunE:  vacationDetailCommand.RunE,
}

var vacationDetailCommand = &cobra.Command{
	Use:   "details",
	Short: "show vacation statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		vacation := application.Conf.GetVacation()

		if vacation.ServiceId == "" {
			return errors.New(textProjectOrServiceNotConfigured)
		}

		serviceId, err := strconv.Atoi(vacation.ServiceId)
		if err != nil {
			return errors.New(textProjectOrServiceNotConfigured)
		}

		entries, err := application.MiteApi.TimeEntries(&domain.TimeEntryQuery{
			At:        entryListFilterAtThisYear,
			ServiceId: domain.NewServiceId(serviceId),
		})
		if err != nil {
			return err
		}

		today := domain.Today()
		var minutesInYear int
		var minutesInPast int
		var minutesInFuture int
		for _, entry := range entries {
			minutesInYear += entry.Minutes.Value()

			if entry.Date.Before(today) {
				minutesInPast += entry.Minutes.Value()
			} else {
				minutesInFuture += entry.Minutes.Value()
			}
		}

		if vacation.Days == "" {
			return errors.New(textProjectOrServiceNotConfigured)
		}
		vacationDays, err := strconv.ParseFloat(vacation.Days, 64)
		if err != nil {
			return errors.New(textProjectOrServiceNotConfigured)
		}

		var daysInYear = domain.MinutesAsDays(minutesInYear, fullVacationDayDuration)
		var daysUnplanned = vacationDays - daysInYear // => user config, if not set explain how

		if vacationDetailsVerbose {
			var daysInPast = domain.MinutesAsDays(minutesInPast, fullVacationDayDuration)
			var daysInFuture = domain.MinutesAsDays(minutesInFuture, fullVacationDayDuration)

			fmt.Printf("Vacation statistics of %d:\n"+
				" - total:     %s days\n"+
				"---------------------\n"+
				" - booked:    %.1f days\n"+
				" - taken:     %.1f days\n"+
				" - planned:   %.1f days\n"+
				" - unplanned: %.1f days\n",
				domain.ThisYear(),
				vacation.Days,
				daysInYear,
				daysInPast,
				daysInFuture,
				daysUnplanned)
		} else {
			fmt.Printf("Vacation statistics of %d:\n"+
				" - booked:    %.1f days\n"+
				" - unplanned: %.1f days\n",
				domain.ThisYear(),
				daysInYear,
				daysUnplanned)
		}

		return nil
	},
}

var vacationCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "creates a vacation entry",
	RunE: func(cmd *cobra.Command, args []string) error {
		vacation := application.Conf.GetVacation()

		if vacation.ProjectId == "" || vacation.ServiceId == "" {
			return errors.New(textProjectOrServiceNotConfigured)
		}

		projectId, err := strconv.Atoi(vacation.ProjectId)
		if err != nil {
			return errors.New(textProjectOrServiceNotConfigured)
		}

		serviceId, err := strconv.Atoi(vacation.ServiceId)
		if err != nil {
			return errors.New(textProjectOrServiceNotConfigured)
		}

		projectIdForVacation := domain.NewProjectId(projectId)
		serviceIdForVacation := domain.NewServiceId(serviceId)

        fromDay := domain.Today()
        if vacationFrom != "" {
          fromDay, err = domain.ParseLocalDate(vacationFrom)
          if err != nil {
              return errors.New(textDateFormatNotCorrect)
          }
        }
      
		minutes := domain.NewMinutesFromHours(fullVacationDayDuration)
		if vacationHalfDay {
			minutes = domain.NewMinutesFromHours(halfVacationDayDuration)
		}

        for i := 0; i < vacationAmount; i++ {
            atDate := fromDay.Add(0, 0, i)
			
            timeEntry := domain.TimeEntryCommand{
				Date:      &atDate,
				Minutes:   &minutes,
				Note:      vacationNote,
				ProjectId: projectIdForVacation,
				ServiceId: serviceIdForVacation,
			}

			_, err := application.MiteApi.CreateTimeEntry(&timeEntry)
			if err != nil {
				return err
			}
        }

        return nil
	},
}

var vacationCountdownCommand = &cobra.Command{
	Use:   "countdown",
	Short: "shows how much time is left to your next vacation",
	RunE: func(cmd *cobra.Command, args []string) error {
		vacation := application.Conf.GetVacation()

		if vacation.ServiceId == "" {
			return errors.New(textProjectOrServiceNotConfigured)
		}

		serviceId, err := strconv.Atoi(vacation.ServiceId)
		if err != nil {
			return errors.New(textProjectOrServiceNotConfigured)
		}

        today := domain.Today()

		entries, err := application.MiteApi.TimeEntries(&domain.TimeEntryQuery{
            From:      &today,
            Direction: "asc",
			ServiceId: domain.NewServiceId(serviceId),
		})
		if err != nil {
			return err
		}

        firstEntry := entries[0]

        unixSeconds := firstEntry.Date.Unix()
        secondsLeft := unixSeconds - today.Unix()
        minutesLeft := secondsLeft / 60
        hoursLeft   := minutesLeft / 60
        daysLeft    := hoursLeft / 24

        withHours   := hoursLeft - (daysLeft * 24)
        withMinutes := minutesLeft - (hoursLeft * 60)
        withSeconds := secondsLeft - (minutesLeft * 60)

        fmt.Printf("Your next vacation starts in:\n"+
            "\n"+
            "\t %dd %dh %dm %ds\n\n",
            daysLeft,
            withHours,
            withMinutes,
            withSeconds)

        if vacationCountdownVerbose {
            fmt.Printf("\t===================\n"+
                "\n"+
                "\t Total: \n"+
                "\t - %d days\n"+
                "\t - %d hours\n"+
                "\t - %d minutes\n"+
                "\t - %d seconds\n\n",
                daysLeft,
                hoursLeft,
                minutesLeft,
                secondsLeft)
	    }

        return nil
	},
}