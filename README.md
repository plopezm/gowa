#GOWA

Gowa is a web admin manager for Go. Currently Gowa is being developed.

#How to use (Example)

```
    /**==================================
		    ADMIN MANAGER MODULE
	 ==================================*/
	GM := gowa.GowaStart("sqlite3", "test.db", 10)

	GM.AddModel(&Driver{})
	GM.AddModel(&Vehicle{})


	for _, r := range GM.GetRoutes(){
		routes = append(routes, r)
	}
	//==================================

	router.PathPrefix("/admin/").Handler(http.StripPrefix("/admin/", http.FileServer(http.Dir(gowa.GowaGetTemplatePath()))))

	log.Fatal(http.ListenAndServe(port, router));

```



#LICENSE

MIT
