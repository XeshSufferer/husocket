package core

/*func ListenWithGracefulShutdown(app *fiber.App, addr string, hs []*Hub) {
	ListenWithGracefulShutdownWithReason(app, addr, hs, "server shutdown")
}*/

/*func ListenWithGracefulShutdownWithReason(app *fiber.App, addr string, hs []*Hub, reason string) {
	go func() {
		if err := app.Listen(addr); err != nil {
			log.Printf("SERVER START ERR: %v", err)
		}
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("SERVER STOP GRACEFUL...")
	wg := sync.WaitGroup{}
	for _, h := range hs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			h.CloseWithReason(reason)
		}()
	}
	wg.Wait()

	time.Sleep(800 * time.Millisecond)
	log.Println("SERVER STOP GRACEFUL...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := app.ShutdownWithContext(ctx)

	if err != nil {
		log.Printf("FIBER STOP ERROR: %v", err)
	}

	log.Println("SHUTTING DOWN FIBER...")

	log.Println("SERVER STOP GRACEFUL! OK.")
}
*/
