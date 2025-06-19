package controller

// simple middleware
// func AuthMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		token := r.Header.Get("Authorization")

// 		if token != "secret-token" {
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}

// 		// Token is valid, call the next handler
// 		next.ServeHTTP(w, r)
// 	})
// }

// func LoggingMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()

// 		// Call the next handler
// 		next.ServeHTTP(w, r)

// 		// Log after the request is processed
// 		log.Printf("%s %s took %s", r.Method, r.RequestURI, time.Since(start))
// 	})
// }
