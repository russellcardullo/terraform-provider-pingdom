package pingdom

// import (
// 	"bytes"
// 	"fmt"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
// 	"github.com/DrFaust92/go-pingdom/solarwinds"
// 	"html/template"
// 	"testing"
// )

// func TestAccUser_basic(t *testing.T) {
// 	email := acctest.RandString(10) + "@foo.com"
// 	resourceName := "pingdom_user.test"
// 	user := solarwinds.User{
// 		Email: email,
// 		Role:  "MEMBER",
// 		Products: []solarwinds.Product{
// 			{
// 				Name: "APPOPTICS",
// 				Role: "MEMBER",
// 			},
// 		},
// 	}

// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckUserDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccUserBasicConfig(user),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckUserExist(resourceName),
// 					resource.TestCheckResourceAttr(resourceName, "role", "MEMBER"),
// 					resource.TestCheckResourceAttr(resourceName, "products.#", "1"),
// 					resource.TestCheckResourceAttr(resourceName, "products.0.name", "APPOPTICS"),
// 					resource.TestCheckResourceAttr(resourceName, "products.0.role", "MEMBER"),
// 				),
// 			},
// 			{
// 				ResourceName:      resourceName,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 		},
// 	})
// }

// func TestAccUser_update(t *testing.T) {
// 	email := acctest.RandString(10) + "@foo.com"
// 	resourceName := "pingdom_user.test"
// 	user := solarwinds.User{
// 		Email: email,
// 		Role:  "MEMBER",
// 		Products: []solarwinds.Product{
// 			{
// 				Name: "APPOPTICS",
// 				Role: "MEMBER",
// 			},
// 		},
// 	}
// 	userUpdate := user
// 	userUpdate.Role = "ADMIN"
// 	userUpdate.Products = append(userUpdate.Products, solarwinds.Product{
// 		Name: "PINGDOM",
// 		Role: "VIEWER",
// 	})

// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckUserDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccUserBasicConfig(user),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckUserExist(resourceName),
// 					resource.TestCheckResourceAttr(resourceName, "role", user.Role),
// 					resource.TestCheckResourceAttr(resourceName, "products.#", fmt.Sprint(len(user.Products))),
// 					resource.TestCheckResourceAttr(resourceName, "products.0.name", user.Products[0].Name),
// 					resource.TestCheckResourceAttr(resourceName, "products.0.role", user.Products[0].Role),
// 				),
// 			},
// 			{
// 				Config: testAccUserBasicConfig(userUpdate),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckUserExist(resourceName),
// 					resource.TestCheckResourceAttr(resourceName, "role", userUpdate.Role),
// 					resource.TestCheckResourceAttr(resourceName, "products.#", fmt.Sprint(len(userUpdate.Products))),
// 					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "products.*", map[string]string{
// 						"name": userUpdate.Products[0].Name,
// 						"role": userUpdate.Products[0].Role,
// 					}),
// 					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "products.*", map[string]string{
// 						"name": userUpdate.Products[1].Name,
// 						"role": userUpdate.Products[1].Role,
// 					}),
// 				),
// 			},
// 			{
// 				ResourceName:      resourceName,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 		},
// 	})
// }

// func testAccCheckUserDestroy(s *terraform.State) error {
// 	client := testAccProvider.Meta().(*Clients).Solarwinds

// 	for _, rs := range s.RootModule().Resources {
// 		if rs.Type != "pingdom_user" {
// 			continue
// 		}

// 		email := rs.Primary.ID
// 		user, err := client.UserService.Retrieve(email)
// 		if err != nil {
// 			return err
// 		}
// 		if user != nil {
// 			return fmt.Errorf("user for resource (%s) still exists", email)
// 		}
// 	}
// 	return nil
// }

// func testAccCheckUserExist(n string) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		rs, ok := s.RootModule().Resources[n]
// 		if !ok {
// 			return fmt.Errorf("not found: %s", n)
// 		}

// 		if rs.Primary.ID == "" {
// 			return fmt.Errorf("no id is set")
// 		}

// 		email := rs.Primary.ID
// 		client := testAccProvider.Meta().(*Clients).Solarwinds
// 		user, err := client.UserService.Retrieve(email)
// 		if err != nil {
// 			return err
// 		}
// 		if user == nil {
// 			return fmt.Errorf("user for resource (%s) not found", email)
// 		}
// 		return nil
// 	}
// }

// func testAccUserBasicConfig(user solarwinds.User) string {
// 	t := template.Must(template.New("basicConfig").Parse(`
// resource "pingdom_user" "test" {
// 	email = "{{.Email}}"
// 	role = "{{.Role}}"
// 	{{range .Products}}
// 	products {
// 		name = "{{.Name}}"
// 		role = "{{.Role}}"
// 	}
// 	{{end}}
// }
// `))
// 	var buf bytes.Buffer
// 	if err := t.Execute(&buf, user); err != nil {
// 		panic(err)
// 	}
// 	result := buf.String()
// 	return result
// }
