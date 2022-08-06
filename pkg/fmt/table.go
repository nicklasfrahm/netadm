package fmt

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/nicklasfrahm/netadm/pkg/nsdp"
)

func Table(output io.Writer, devices []nsdp.Device, columns []string) {
	// Normalize column names.
	for i, column := range columns {
		columns[i] = strings.ToLower(column)
	}

	// Create table with tabwriter.
	w := tabwriter.NewWriter(output, 0, 0, 4, ' ', tabwriter.TabIndent)

	// Fetch table columns from desired keys.
	for _, column := range columns {
		// Print column header.
		rt := nsdp.RecordTypeByName[column]
		fmt.Fprintf(w, "%s\t", strings.ToUpper(rt.Name))
	}
	fmt.Fprintln(w)

	// Print requested keys for each device. Note that we
	// unmarshal the message into a Device because this
	// allows it to easily group messages of the same type.
	for _, device := range devices {
		// Print the desired columns.
		for _, column := range columns {
			// Fetch field from device.
			name := nsdp.RecordTypeByName[column].Name
			field := reflect.ValueOf(device).FieldByName(name)
			if field.IsValid() {
				if field.Kind() == reflect.String && field.IsZero() {
					fmt.Fprintf(w, "<nil>\t")
				} else {
					fmt.Fprintf(w, "%v\t", field)
				}
			} else {
				// This happens if the field is a known message
				// type but not defined inside the Device struct.
				// Make sure to add the according field to the
				// Device struct to prevent this.
				fmt.Fprintf(w, "<invalid>\t")
			}
		}

		fmt.Fprintln(w)
	}

	w.Flush()
}
