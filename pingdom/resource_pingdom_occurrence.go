package pingdom

import (
	"context"
	"fmt"
	"github.com/DrFaust92/go-pingdom/pingdom"
	"github.com/DrFaust92/go-pingdom/solarwinds"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"time"
)

func resourcePingdomOccurrences() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePingdomOccurrencesCreate,
		ReadContext:   resourcePingdomOccurrencesRead,
		UpdateContext: resourcePingdomOccurrencesUpdate,
		DeleteContext: resourcePingdomOccurrencesDelete,
		Schema: map[string]*schema.Schema{
			"maintenance_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"effective_from": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"effective_to": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"from": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"to": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
		},
	}
}

type OccurrenceGroup pingdom.ListOccurrenceQuery

func timeParse(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}

func timeFormat(unixTime int64) string {
	return time.Unix(unixTime, 0).Format(time.RFC3339)
}

func getTime(attr string, d *schema.ResourceData) (int64, bool, error) {
	v, ok := d.GetOk(attr)
	if ok {
		t, err := timeParse(v.(string))
		if err != nil {
			return 0, false, err
		}
		return t.Unix(), true, nil
	}
	return 0, false, nil
}

func NewOccurrenceGroupWithResourceData(d *schema.ResourceData) (*OccurrenceGroup, error) {
	q := OccurrenceGroup{}

	// required
	if v, ok := d.GetOk("maintenance_id"); ok {
		q.MaintenanceId = int64(v.(int))
	}

	if v, ok, err := getTime("effective_from", d); err != nil {
		return nil, err
	} else if ok {
		q.From = v
	}

	if v, ok, err := getTime("effective_to", d); err != nil {
		return nil, err
	} else if ok {
		q.To = v
	}

	return &q, nil
}

// ID returns unique id for an OccurrenceGroup. An OccurrenceGroup is essentially a query against Maintenance Occurrence.
// The result of query can overlap, so there is no unique resource id for queries on the Pingdom side.
func (g *OccurrenceGroup) ID() string {
	return solarwinds.RandString(32)
}

func (g *OccurrenceGroup) List(client *pingdom.Client) ([]pingdom.Occurrence, error) {
	return client.Occurrences.List(pingdom.ListOccurrenceQuery(*g))
}

func (g *OccurrenceGroup) Populate(client *pingdom.Client, d *schema.ResourceData) error {
	sample, size, err := g.Sample(client)
	if err != nil {
		return err
	}
	for k, v := range map[string]interface{}{
		"from":           timeFormat(sample.From),
		"to":             timeFormat(sample.To),
		"effective_from": timeFormat(g.From),
		"effective_to":   timeFormat(g.To),
		"maintenance_id": g.MaintenanceId,
		"size":           size,
	} {
		if err = d.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (g *OccurrenceGroup) Sample(client *pingdom.Client) (*pingdom.Occurrence, int, error) {
	occurrences, err := g.List(client)
	if err != nil {
		return nil, 0, err
	}
	if len(occurrences) == 0 {
		return nil, 0, fmt.Errorf("there are no occurrences matching query: %#v", g)
	}
	return &occurrences[0], len(occurrences), nil
}

func (g *OccurrenceGroup) Size(client *pingdom.Client) (int, error) {
	occurrences, err := g.List(client)
	if err != nil {
		return 0, err
	}

	return len(occurrences), nil
}

func (g *OccurrenceGroup) MustExists(client *pingdom.Client) error {
	size, err := g.Size(client)
	if err != nil {
		return err
	}
	if size == 0 {
		return fmt.Errorf("there are no occurrences matching query: %#v", g)
	}
	return nil
}

func (g *OccurrenceGroup) Update(client *pingdom.Client, from int64, to int64) error {
	occurrenceUpdate := pingdom.Occurrence{
		From: from,
		To:   to,
	}
	return g.groupOp(client, func(occurrence pingdom.Occurrence) (interface{}, error) {
		return client.Occurrences.Update(occurrence.Id, occurrenceUpdate)
	})
}

func (g *OccurrenceGroup) Delete(client *pingdom.Client) error {
	return g.groupOp(client, func(occurrence pingdom.Occurrence) (interface{}, error) {
		return client.Occurrences.Delete(occurrence.Id)
	})
}

func (g *OccurrenceGroup) groupOp(client *pingdom.Client, op func(occurrence pingdom.Occurrence) (interface{}, error)) error {
	occurrences, err := client.Occurrences.List(pingdom.ListOccurrenceQuery(*g))
	if err != nil {
		return err
	}

	cancelChan := make(chan bool)
	errChan := make(chan error, len(occurrences))
	for _, occurrence := range occurrences {
		go func(occurrence pingdom.Occurrence) {
			select {
			case <-cancelChan:
				return
			default:
				_, err := op(occurrence)
				errChan <- err
			}
		}(occurrence)
	}

	expectTotal := len(occurrences)
	count := 0
	for err := range errChan {
		if err != nil {
			close(cancelChan)
			return err
		} else {
			count++
		}
		if expectTotal == count {
			break
		}
	}
	return nil
}

func resourcePingdomOccurrencesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).Pingdom

	g, err := NewOccurrenceGroupWithResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	var from, to int64
	if v, ok, err := getTime("from", d); err != nil {
		return diag.FromErr(err)
	} else if ok {
		from = v
	}
	if v, ok, err := getTime("to", d); err != nil {
		return diag.FromErr(err)
	} else if ok {
		to = v
	}

	if (from == 0 && to != 0) || (from != 0 && to == 0) {
		return diag.Errorf("'from' and 'to' must be provided at the same time, current values are from: %d, to: %d", from, to)
	}

	log.Printf("[DEBUG] Retrieve occurrences with query: %#v", g)
	if err := g.Populate(client, d); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(g.ID())

	var newFrom, newTo int64
	if v, ok, err := getTime("from", d); err != nil {
		return diag.FromErr(err)
	} else if ok {
		newFrom = v
	}
	if v, ok, err := getTime("to", d); err != nil {
		return diag.FromErr(err)
	} else if ok {
		newTo = v
	}

	if from != 0 && to != 0 && (from != newFrom || to != newTo) {
		log.Printf("User specifies new 'from' and 'to' upon creation, from: %d, to: %d", from, to)
		if err := g.Update(client, from, to); err != nil {
			return diag.FromErr(err)
		}
		if err := g.Populate(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourcePingdomOccurrencesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).Pingdom

	g, err := NewOccurrenceGroupWithResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Retrieve occurrences with query: %#v", g)
	if err := g.Populate(client, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePingdomOccurrencesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).Pingdom

	g, err := NewOccurrenceGroupWithResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	var updated bool

	if d.HasChanges("effective_from") || d.HasChanges("effective_to") {
		log.Printf("[DEBUG] Retrieve occurrences with query: %#v", g)
		if err := g.Populate(client, d); err != nil {
			return diag.FromErr(err)
		}
		updated = true
	}

	if d.HasChanges("from") || d.HasChanges("to") {
		var from, to int64
		if v, ok, err := getTime("from", d); err != nil {
			return diag.FromErr(err)
		} else if ok {
			from = v
		}
		if v, ok, err := getTime("to", d); err != nil {
			return diag.FromErr(err)
		} else if ok {
			to = v
		}

		if from == 0 || to == 0 {
			return diag.Errorf("'from' and 'to' must be provided at the same time, current values are from: %d, to: %d", from, to)
		}

		log.Printf("[DEBUG] Occurrence update from: %d, to: %d", from, to)

		if err := g.Update(client, from, to); err != nil {
			return diag.FromErr(err)
		}
		updated = true
	}

	if updated {
		if err := g.Populate(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourcePingdomOccurrencesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).Pingdom

	occurrence, err := NewOccurrenceGroupWithResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = occurrence.Delete(client)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
