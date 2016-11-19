#GOWA

Gowa is a web admin manager for Go. Currently Gowa is being developed.

#How to use (Example)

```
    /**==================================
    		ADMIN MANAGER MODULE
     ==================================*/
        //Starting gowa with SQLite3 driver
    	GM := gowa.GowaStart("sqlite3", "test.db", 10)

        //Adding modules to the admin manager
    	GM.AddModel(&User{})
    	GM.AddModel(&Company{})
    	GM.AddModel(&Driver{})
    	GM.AddModel(&Vehicle{})

    	//Generate first user
    	gowa.GowaCreateAdminUser("pabloplm@gmail.com", "example")

        //The router is a *mux.Router
    	gowa.GowaAddRoutes(router);

    	//Then we add the template in whatever we need, in this case, the path will be /admin
    	router.PathPrefix("/admin/").Handler(http.StripPrefix("/admin/", http.FileServer(http.Dir(gowa.GowaGetTemplatePath()))))


```



#LICENSE

MIT
