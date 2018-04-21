/*
Package pingdom provides a client interface to the Pingdom API.  This currently only
supports working with basic HTTP and ping checks.

Construct a new Pingdom client:

	client := pingdom.NewClient("pingdom_username", "pingdom_password", "pingdom_api_key")

Using a Pingdom client, you can access supported services.

CheckService

This service manages pingdom Checks which are represented by the `Check` struct.
When creating or updating Checks you must specify at a minimum the `Name`, `Hostname`
and `Resolution`.  Other fields are optional but if not set will be given the zero
values for the underlying type.

More information on Checks from Pingdom: https://www.pingdom.com/features/api/documentation/#ResourceChecks

Get a list of all checks:

	checks, err := client.Checks.List()
	fmt.Println("Checks:", checks) // [{ID Name} ...]

Create a new HTTP check:

	newCheck := pingdom.Check{Name: "Test Check", Hostname: "example.com", Resolution: 5}
	check, err := client.Checks.Create(&newCheck)
	fmt.Println("Created check:", check) // {ID, Name}

Create a new HTTP check with alerts for specified users:

	newCheck := pingdom.Check{Name: "Test Check", Hostname: "example.com", Resolution: 5, UserIds: []int{12345}}
	check, err := client.Checks.Create(&newCheck)
	fmt.Println("Created check:", check) // {ID, Name}

Create a new Ping check:

	newCheck := pingdom.PingCheck{Name: "Test Check", Hostname: "example.com", Resolution: 5}
	check, err := client.Checks.Create(&newCheck)
	fmt.Println("Created check:", check) // {ID, Name}

Get details for a specific check:

	checkDetails, err := client.Checks.Read(12345)

Update a check:

	updatedCheck := pingdom.Check{Name: "Updated Check", Hostname: "example2.com", Resolution: 5}
	msg, err := client.Checks.Update(12345, &updatedCheck)

Delete a check:

	msg, err := client.Checks.Delete(12345)

*/
package pingdom
