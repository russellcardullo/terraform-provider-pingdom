package pingdom

import (
	"bytes"
	"fmt"
	"github.com/DrFaust92/go-pingdom/solarwinds"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strconv"
	"testing"
	"text/template"
	"time"
)

func TestAccOccurrence_basic(t *testing.T) {
	occurrenceNum := 3
	maintenance := getMaintenance(time.Duration(occurrenceNum))
	resourceName := "pingdom_occurrence.test"

	var from, to string
	from = maintenance["From"]
	if v, err := timeParse(maintenance["To"]); err != nil {
		t.Fatal(err)
	} else {
		to = timeFormat(v.Add(1 * time.Hour).Unix())
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOccurrenceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOccurrenceBasicConfig(maintenance, from, to),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "from", from),
					resource.TestCheckResourceAttr(resourceName, "to", to),
					resource.TestCheckResourceAttr(resourceName, "size", strconv.Itoa(occurrenceNum+1)),
				),
			},
		},
	})
}

func TestAccOccurrence_update(t *testing.T) {
	occurrenceNum := 3
	maintenance := getMaintenance(time.Duration(occurrenceNum))
	resourceName := "pingdom_occurrence.test"

	var from, to string
	from = maintenance["From"]
	if v, err := timeParse(maintenance["To"]); err != nil {
		t.Fatal(err)
	} else {
		to = timeFormat(v.Add(1 * time.Hour).Unix())
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOccurrenceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOccurrenceBasicConfig(maintenance, "", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "from", maintenance["From"]),
					resource.TestCheckResourceAttr(resourceName, "to", maintenance["To"]),
					resource.TestCheckResourceAttr(resourceName, "size", strconv.Itoa(occurrenceNum+1)),
				),
			},
			{
				Config: testAccOccurrenceBasicConfig(maintenance, from, to),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "from", from),
					resource.TestCheckResourceAttr(resourceName, "to", to),
					resource.TestCheckResourceAttr(resourceName, "size", strconv.Itoa(occurrenceNum+1)),
				),
			},
		},
	})
}

func testAccCheckOccurrenceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Clients).Pingdom

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingdom_occurrence" {
			continue
		}

		g := OccurrenceGroup{}
		v, err := strconv.ParseInt(rs.Primary.Attributes["maintenance_id"], 10, 64)
		if err != nil {
			return err
		}
		g.MaintenanceId = v

		t, err := timeParse(rs.Primary.Attributes["effective_from"])
		if err != nil {
			return err
		}
		g.From = t.Unix()

		t, err = timeParse(rs.Primary.Attributes["effective_to"])
		if err != nil {
			return err
		}
		g.To = t.Unix()

		if size, err := g.Size(client); err != nil {
			return err
		} else if size != 0 {
			return fmt.Errorf("the occurrence has not been deleted, %d left", size)
		}
	}
	return nil
}

func getMaintenance(occurrenceNum time.Duration) map[string]string {
	now := time.Now()
	from := now.Add(1 * time.Hour)
	to := from.Add(1 * time.Hour)
	return map[string]string{
		"Description": "terraform resource test - " + solarwinds.RandString(10),
		"From":        timeFormat(from.Unix()),
		"To":          timeFormat(to.Unix()),
		"EffectiveTo": timeFormat(to.Add(occurrenceNum * 24 * time.Hour).Unix()),
	}
}

func testAccOccurrenceBasicConfig(maintenance map[string]string, from string, to string) string {
	t := template.Must(template.New("basicConfig").Parse(`
{{with .maintenance}}
resource "pingdom_maintenance" "test" {
	description = "{{.Description}}"
	from = "{{.From}}"
	to = "{{.To}}"
	recurrencetype = "day"
	repeatevery = 1
	effectiveto = "{{.EffectiveTo}}"
}
{{end}}

resource "pingdom_occurrence" "test" {
	maintenance_id = pingdom_maintenance.test.id
	effective_from = pingdom_maintenance.test.from
	effective_to = pingdom_maintenance.test.effectiveto
	{{if .from}}
	from = "{{.from}}"
	{{end}}
	{{if .to}}
	to = "{{.to}}"
	{{end}}
}
`))
	var buf bytes.Buffer
	if err := t.Execute(&buf, map[string]interface{}{
		"maintenance": maintenance,
		"from":        from,
		"to":          to,
	}); err != nil {
		panic(err)
	}
	result := buf.String()
	return result
}
