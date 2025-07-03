package mysqlprofiler

import (
	"database/sql"
	"fmt"
	"math"
	"strings"

	"github.com/eiiches/mysql-protobuf-functions/internal/moremaps"
)

func RunProfile(db *sql.DB) {
	stmt, err := db.Prepare("SELECT event_id, nesting_event_id, event_name, object_type, object_name, sql_text, timer_wait, timer_start, timer_end FROM performance_schema.events_statements_history_long")
	if err != nil {
		panic(err)
	}
	defer func() {
		stmt.Close()
	}()

	rows, err := stmt.Query()
	if err != nil {
		panic(err)
	}

	type Event struct {
		EventID        uint64
		NestingEventID *uint64
		EventName      string
		ObjectType     *string
		ObjectName     *string
		SqlText        *string
		TimerWait      int64 // in picoseconds
		TimerStart     int64
		TimerEnd       int64
	}

	events := map[uint64]*Event{}

	for rows.Next() {
		event := &Event{}
		if err := rows.Scan(&event.EventID, &event.NestingEventID, &event.EventName, &event.ObjectType, &event.ObjectName, &event.SqlText, &event.TimerWait, &event.TimerStart, &event.TimerEnd); err != nil {
			panic(err)
		}
		events[event.EventID] = event
		if len(events)%100000 == 0 {
			fmt.Printf("Loaded %d events\n", len(events))
		}
	}

	collapsed := map[string]int64{}

	for _, leafEvent := range events {
		stack := []*Event{leafEvent}

		event := leafEvent
		for event.NestingEventID != nil {
			if parent, ok := events[*event.NestingEventID]; ok {
				stack = append(stack, parent)
				event = parent
			} else {
				break // No parent found, exit the loop
			}
		}

		parentStackText := ""
		stackText := strings.Builder{}
		sep := ""
		for i := len(stack) - 1; i >= 0; i-- {
			stackText.WriteString(sep)
			// if i == 0 {
			if stack[i].SqlText != nil {
				stackText.WriteString(strings.ReplaceAll(strings.ReplaceAll(*stack[i].SqlText, "\n", " "), ";", "$"))
			} else {
				stackText.WriteString(stack[i].EventName)
			}
			stackText.WriteString(" in ")
			if stack[i].ObjectType != nil {
				stackText.WriteString(*stack[i].ObjectType)
				stackText.WriteString(":")
				stackText.WriteString(*stack[i].ObjectName)
			} else {
				stackText.WriteString("main")
			}

			if i == 1 {
				parentStackText = stackText.String()
			}
			sep = ";"
		}

		collapsed[parentStackText] -= leafEvent.TimerWait
		collapsed[stackText.String()] += leafEvent.TimerWait
	}

	for stack, totalWait := range moremaps.SortedEntries(collapsed) {
		if stack == "" {
			continue // Skip empty stacks
		}
		fmt.Printf("%s %d\n", stack, int64(math.Max(float64(totalWait/1000000), 1))) // picoseconds -> microseconds
	}
}
