package tsbuild

import (
	"time"

	"github.com/delaneyj/toolbelt"
)

type EnumValueDefinition struct {
	Name        toolbelt.CasedString
	Value       string
	Description string
}

type EnumDefinition struct {
	Name         toolbelt.CasedString
	Values       []*EnumValueDefinition
	DefaultIndex int
	Default      *EnumValueDefinition
}

type ConstTemplateData struct {
	Version                   string
	DatastarKey               string
	VersionClientByteSize     int
	VersionClientByteSizeGzip int
	DefaultSettleDuration     time.Duration
	DefaultSSERetryDuration   time.Duration
	DefaultUseViewTransitions bool
	Enums                     []*EnumDefinition
}

var ConstsData = &ConstTemplateData{
	DatastarKey:               "datastar",
	DefaultSettleDuration:     300 * time.Millisecond,
	DefaultSSERetryDuration:   1 * time.Second,
	DefaultUseViewTransitions: false,
	Enums: []*EnumDefinition{
		{
			Name:         toolbelt.ToCasedString("FragmentMergeMode"),
			DefaultIndex: 0,
			Values: []*EnumValueDefinition{
				{
					Value:       "morph",
					Description: "Morphs the fragment into the existing element using idiomorph.",
				},
				{
					Value:       "inner",
					Description: "Replaces the inner HTML of the existing element.",
				},
				{
					Value:       "outer",
					Description: "Replaces the outer HTML of the existing element.",
				},
				{
					Value:       "prepend",
					Description: "Prepends the fragment to the existing element.",
				},
				{
					Value:       "append",
					Description: "Appends the fragment to the existing element.",
				},
				{
					Value:       "before",
					Description: "Inserts the fragment before the existing element.",
				},
				{
					Value:       "after",
					Description: "Inserts the fragment after the existing element.",
				},
				{
					Value:       "upsertAttributes",
					Description: "Upserts the attributes of the existing element.",
				},
			},
		},

		{
			Name:         toolbelt.ToCasedString("EventType"),
			DefaultIndex: -1,
			Values: []*EnumValueDefinition{
				{
					Name:        toolbelt.ToCasedString("Fragment"),
					Value:       "datastar-fragment",
					Description: "An event dealing with HTML fragments",
				},
				{
					Name:        toolbelt.ToCasedString("Signal"),
					Value:       "datastar-signal",
					Description: "An event dealing with fine grain signals",
				},
				{
					Name:        toolbelt.ToCasedString("Remove"),
					Value:       "datastar-remove",
					Description: "An event dealing with removing elements or signals",
				},
				{
					Name:        toolbelt.ToCasedString("Redirect"),
					Value:       "datastar-redirect",
					Description: "An event dealing with redirecting the browser",
				},
				{
					Name:        toolbelt.ToCasedString("Console"),
					Value:       "datastar-console",
					Description: "An event dealing with console messages",
				},
			},
		},

		{
			Name:         toolbelt.ToCasedString("ConsoleMode"),
			DefaultIndex: -1,
			Values: []*EnumValueDefinition{
				{
					Value:       "assert",
					Description: "Writes an error message to the console if the assertion is false. If the assertion is true, nothing happens.",
				},
				{
					Value:       "clear",
					Description: "Clears the console if the console allows it. A graphical console, like those running on browsers, will allow it; a console displaying on the terminal, like the one running on Node, will not support it, and will have no effect (and no error).",
				},
				{
					Value:       "count",
					Description: "Logs the number of times that this particular call to count() has been called.",
				},
				{
					Value:       "countReset",
					Description: "Resets counter used with console.count().",
				},
				{
					Value:       "debug",
					Description: "Outputs a message to the console at the 'debug' log level. The message is only displayed to the user if the console is configured to display debug output. In most cases, the log level is configured within the console UI. This log level might correspond to the Debug or Verbose log level.",
				},
				{
					Value:       "dir",
					Description: "Displays a list of the properties of the specified JavaScript object. In browser consoles, the output is presented as a hierarchical listing with disclosure triangles that let you see the contents of child objects. Unlike other logging methods, console.dir() does not attempt to pretty-print the object. For example, if you pass a DOM element to console.dir(), it will not be displayed like in the element inspector, but will instead show a list of properties.",
				},
				{
					Value:       "dirxml",
					Description: "Displays an interactive tree of the descendant elements of the specified XML/HTML element. If it is not possible to display as an element the JavaScript Object view is shown instead. The output is presented as a hierarchical listing of expandable nodes that let you see the contents of child nodes.",
				},
				{
					Value:       "error",
					Description: "Outputs a message to the console at the 'error' log level.",
				},
				{
					Value:       "group",
					Description: "Creates a new inline group in the Web console log, causing any subsequent console messages to be indented by an additional level, until console.groupEnd() is called.",
				},
				{
					Value:       "groupCollapsed",
					Description: "Creates a new inline group in the console. Unlike console.group(), however, the new group is created collapsed. The user will need to use the disclosure button next to it to expand it, revealing the entries created in the group. Call console.groupEnd() to back out to the parent group.",
				},
				{
					Value:       "groupEnd",
					Description: "Exits the current inline group in the console.",
				},
				{
					Value:       "info",
					Description: "Outputs a message to the console at the 'info' log level.",
				},
				{
					Value:       "log",
					Description: "Outputs a message to the console.",
				},
				{
					Value:       "table",
					Description: "Displays tabular data as a table.",
				},
				{
					Value:       "time",
					Description: "Starts a timer you can use to track how long an operation takes. You give each timer a unique name, and may have up to 10,000 timers running on a given page. When you call console.timeEnd() with the same name, the browser will output the time, in milliseconds, that elapsed since the timer was started.",
				},
				{
					Value:       "timeEnd",
					Description: "The console.timeEnd() static method stops a timer that was previously started by calling console.time().",
				},
				{
					Value:       "timeLog",
					Description: "Logs the current value of a timer that was previously started by calling console.time().",
				},
				{
					Value:       "trace",
					Description: "Outputs a stack trace to the console.",
				},
				{
					Value:       "warn",
					Description: "Outputs a warning message to the console at the 'warning' log level.",
				},
			},
		},
	},
}
