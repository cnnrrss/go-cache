/*

Package go-cache is an in-memory key-value store where keys expire after a period
of time.

The Cache interface provides Get and Set methods. The cache is goroutine safe.

Values are stored as type interface{}. You will need to use a type assertion
after retrieval to use the typed value.

Example:

    c := New(10 * time.Second)
    c.Set("key_name", 1)
    value, found := c.Get("key_name")
    if found && value.(int) == 1 {
        fmt.Println("success")
    }
*/
package cache

