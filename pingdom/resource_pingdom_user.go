package pingdom

// import (
// 	"context"
// 	"errors"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// 	"github.com/DrFaust92/go-pingdom/solarwinds"
// 	"log"
// 	"time"
// )

// const (
// 	DeleteUserRetryTimeout = 1 * time.Minute
// )

// func resourceSolarwindsUser() *schema.Resource {
// 	return &schema.Resource{
// 		CreateContext: resourceSolarwindsUserCreate,
// 		ReadContext:   resourceSolarwindsUserRead,
// 		UpdateContext: resourceSolarwindsUserUpdate,
// 		DeleteContext: resourceSolarwindsUserDelete,
// 		Importer: &schema.ResourceImporter{
// 			StateContext: schema.ImportStatePassthroughContext,
// 		},
// 		Schema: map[string]*schema.Schema{
// 			"email": {
// 				Type:     schema.TypeString,
// 				Required: true,
// 				ForceNew: true,
// 			},
// 			"role": {
// 				Type:     schema.TypeString,
// 				Required: true,
// 				ForceNew: false,
// 			},
// 			"products": {
// 				Type:     schema.TypeSet,
// 				Optional: true,
// 				ForceNew: false,
// 				Elem: &schema.Resource{
// 					Schema: map[string]*schema.Schema{
// 						"name": {
// 							Type:     schema.TypeString,
// 							Required: true,
// 						},
// 						"role": {
// 							Type:     schema.TypeString,
// 							Required: true,
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// }

// func userFromResource(d *schema.ResourceData) (*solarwinds.User, error) {
// 	user := solarwinds.User{}

// 	// required
// 	if v, ok := d.GetOk("email"); ok {
// 		user.Email = v.(string)
// 	}

// 	if v, ok := d.GetOk("role"); ok {
// 		user.Role = v.(string)
// 	}

// 	if v, ok := d.GetOk("products"); ok {
// 		interfaceSlice := v.(*schema.Set).List()
// 		user.Products = expandUserProducts(interfaceSlice)
// 	}

// 	return &user, nil
// }

// func expandUserProducts(l []interface{}) []solarwinds.Product {
// 	if len(l) == 0 || l[0] == nil {
// 		return nil
// 	}

// 	m := make([]solarwinds.Product, 0, len(l))
// 	for _, tfMapRaw := range l {
// 		tfMap, ok := tfMapRaw.(map[string]interface{})
// 		if !ok {
// 			continue
// 		}
// 		product := solarwinds.Product{}
// 		if name, ok := tfMap["name"].(string); ok && name != "" {
// 			product.Name = name
// 		}
// 		if role, ok := tfMap["role"].(string); ok && role != "" {
// 			product.Role = role
// 		}
// 		m = append(m, product)
// 	}

// 	return m
// }

// func resourceSolarwindsUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	client := meta.(*Clients).Solarwinds

// 	user, err := userFromResource(d)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}

// 	log.Printf("[DEBUG] User create configuration: %#v", d.Get("email"))
// 	err = client.UserService.Create(*user)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}

// 	d.SetId(user.Email)
// 	return resourceSolarwindsUserRead(ctx, d, meta)
// }

// func resourceSolarwindsUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	client := meta.(*Clients).Solarwinds

// 	email := d.Id()
// 	user, err := client.UserService.Retrieve(email)
// 	if err != nil {
// 		return diag.Errorf("error retrieving user with email %v", email)
// 	}
// 	if user == nil {
// 		d.SetId("")
// 		return nil
// 	}

// 	for k, v := range map[string]interface{}{
// 		"email":    user.Email,
// 		"role":     user.Role,
// 		"products": flattenUserProducts(user.Products),
// 	} {
// 		if err := d.Set(k, v); err != nil {
// 			return diag.FromErr(err)
// 		}
// 	}

// 	return nil
// }

// func flattenUserProducts(l []solarwinds.Product) []interface{} {
// 	if l == nil {
// 		return []interface{}{}
// 	}
// 	sets := make([]interface{}, 0, len(l))
// 	for _, item := range l {
// 		tfMap := map[string]interface{}{
// 			"name": item.Name,
// 			"role": item.Role,
// 		}
// 		sets = append(sets, tfMap)
// 	}

// 	return sets
// }

// func resourceSolarwindsUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	client := meta.(*Clients).Solarwinds
// 	if d.HasChanges("role", "products") {
// 		user, err := userFromResource(d)
// 		if err != nil {
// 			return diag.FromErr(err)
// 		}

// 		log.Printf("[DEBUG] User update configuration: %#v", user)

// 		if err = client.UserService.Update(*user); err != nil {
// 			return diag.Errorf("Error updating user: %s", err)
// 		}
// 	}

// 	return resourceSolarwindsUserRead(ctx, d, meta)
// }

// func resourceSolarwindsUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	client := meta.(*Clients).Solarwinds

// 	id := d.Id()
// 	err := resource.RetryContext(ctx, DeleteUserRetryTimeout, func() *resource.RetryError {
// 		if err := client.UserService.Delete(id); err != nil {
// 			var clientErr *solarwinds.ClientError
// 			ok := errors.As(err, &clientErr)
// 			if ok && solarwinds.ErrCodeDeleteActiveUserException == clientErr.StatusCode {
// 				return resource.NonRetryableError(err)
// 			}
// 			return resource.RetryableError(err)
// 		}
// 		return nil
// 	})

// 	if err != nil {
// 		return diag.Errorf("error deleting user: %s", err)
// 	}

// 	return nil
// }
