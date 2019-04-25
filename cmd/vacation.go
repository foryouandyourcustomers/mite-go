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
	textProjectOrServiceNotConfigured = "please set both the vacation project AND service id (either via arguments or config)"
)

var (
	vacationDetailsverbose bool
	vacationNote           string
	vacationHalfDay        bool
	vacationFrom           string
	vacationTo             string
)

func init() {
	vacationDetailCommand.Flags().BoolVarP(&vacationDetailsverbose, "verbose", "v", false, "verbose output")
	vacationCommand.AddCommand(vacationDetailCommand)

	vacationCreateCommand.Flags().StringVarP(&vacationNote, "note", "n", "", "a note describing your vacation")
	vacationCreateCommand.Flags().BoolVarP(&vacationHalfDay, "halfday", "d", false, "If set vacation is entered as half a day")
	vacationCreateCommand.Flags().StringVarP(&vacationFrom, "from", "f", "", "create vacation starting at date (in YYYY-MM-DD format)")
	vacationCreateCommand.Flags().StringVarP(&vacationTo, "to", "t", "", "create vacation ending at date (in YYYY-MM-DD format)")
	vacationCommand.AddCommand(vacationCreateCommand)

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
		vacationActivity := application.Conf.GetVacation()

		if vacationActivity.ServiceId == "" {
			return errors.New(textProjectOrServiceNotConfigured)
		}

		serviceId, err := strconv.Atoi(vacationActivity.ServiceId)
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

		var daysInYear = domain.MinutesAsDays(minutesInYear, fullVacationDayDuration)
		var daysInPast = domain.MinutesAsDays(minutesInPast, fullVacationDayDuration)
		var daysInFuture = domain.MinutesAsDays(minutesInFuture, fullVacationDayDuration)
		var daysUnplanned = 28 - daysInYear // => user config, if not set explain how

		if vacationDetailsverbose {
			fmt.Printf("Vacation statistics of %d:\n"+
				" - total:     %d days\n"+
				"---------------------\n"+
				" - booked:    %.1f days\n"+
				" - taken:     %.1f days\n"+
				" - planned:   %.1f days\n"+
				" - unplanned: %.1f days\n",
				domain.ThisYear(),
				28,
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
	Short: "creates a vacation entry (WIP: currently this command creates a vacation day only for today)",
	RunE: func(cmd *cobra.Command, args []string) error {
		vacationActivity := application.Conf.GetVacation()

		if vacationActivity.ProjectId == "" || vacationActivity.ServiceId == "" {
			return errors.New(textProjectOrServiceNotConfigured)
		}

		projectId, err := strconv.Atoi(vacationActivity.ProjectId)
		if err != nil {
			return errors.New(textProjectOrServiceNotConfigured)
		}

		serviceId, err := strconv.Atoi(vacationActivity.ServiceId)
		if err != nil {
			return errors.New(textProjectOrServiceNotConfigured)
		}

		projectIdForVacation := domain.NewProjectId(projectId)
		serviceIdForVacation := domain.NewServiceId(serviceId)
		today := domain.Today()

		minutes := domain.NewMinutesFromHours(fullVacationDayDuration)
		if vacationHalfDay {
			minutes = domain.NewMinutesFromHours(halfVacationDayDuration)
		}

		var dates []domain.LocalDate
		dates = append(dates, today)

		for _, date := range dates {
			timeEntry := domain.TimeEntryCommand{
				Date:      &date,
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
