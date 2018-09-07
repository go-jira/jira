package jiradata
import "fmt"

type Ticket struct {
	Title       string
	Description string
}

func NewTicket(title, description string) Ticket {
	return Ticket{Title: title, Description: description}
}

func (ticket Ticket) String() string {
	return fmt.Sprintf("{\ntitle: %s\ndescription: %s\n}", ticket.Title, ticket.Description)
}

func (ticket Ticket) stringIndent() string {
	return fmt.Sprintf("\t\t{\n\t\t\ttitle: %s\n\t\t\tdescription: %s\n\t\t}", ticket.Title, ticket.Description)
}

type Epic struct {
	Title       string
	Summary     string
	Description string
	Tickets     []Ticket
}

func NewEpic(title, summary string, description string) Epic {
	return Epic{Title: title, Summary: summary, Description: description}
}

func (epic Epic) String() string {
	var ticketsString string
	for _, t := range epic.Tickets {
		ticketsString += t.stringIndent() + ",\n"
	}

	if ticketsString != "" {
		ticketsString = ticketsString[:len(ticketsString) - 2]
	}

	return fmt.Sprintf(
		"{\n\ttitle: %s\n\tsummary: %s\n\tdescription: %s\n\ttickets: %s\n}",
		epic.Title, epic.Summary, epic.Description, ticketsString)
}

func (epic *Epic) AddTicket(ticket Ticket) {
	epic.Tickets = append(epic.Tickets, ticket)
}

type Plan struct {
	Epics []Epic
}

func (plan Plan) String() string {
	var epicsString string
	for _, e := range plan.Epics {
		epicsString += e.String() + ",\n"
	}

	if epicsString != "" {
		epicsString = epicsString[:len(epicsString) - 2]
	}

	return fmt.Sprintf("{\nepics: %s\n}", epicsString)
}

func (plan *Plan) AddEpic(epic Epic) {
	plan.Epics = append(plan.Epics, epic)
}