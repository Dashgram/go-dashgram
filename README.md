# Go Dashgram SDK

A Go SDK for integrating with the [Dashgram](https://dashgram.io) analytics platform. This library provides both synchronous and asynchronous methods for tracking Telegram bot events and user invitations.

## Features

- **Synchronous & Asynchronous Tracking**: Track events with both blocking and non-blocking methods
- **Flexible Configuration**: Customize API URLs, HTTP clients, and worker pools
- **Context Support**: Full context.Context support for cancellation and timeouts
- **Error Handling**: Comprehensive error types for different failure scenarios
- **Easy Integration**: Simple integration with popular Telegram bot libraries

## Installation

```bash
go get github.com/Dashgram/go-dashgram
```

## Quick Start

### Basic Usage

```go
package main

import (
    "log"
    "github.com/Dashgram/go-dashgram"
)

func main() {
    client := dashgram.New(12345, "your-access-key")
    defer client.Close()

    event := map[string]interface{}{
        "user_id": 123456789,
        "action":  "message_sent",
        "chat_id": 987654321,
    }

    if err := client.TrackEvent(event); err != nil {
        log.Printf("Failed to track event: %v", err)
    }
}
```

### With Custom Configuration

```go
package main

import (
    "log"
    "time"
    "net/http"
    "github.com/Dashgram/go-dashgram"
)

func main() {
    // Create custom HTTP client
    httpClient := &http.Client{
        Timeout: 10 * time.Second,
    }

    // Create client with custom options
    client := dashgram.New(12345, "your-access-key",
        dashgram.WithAPIURL("https://custom-api.dashgram.io/v1"),
        dashgram.WithOrigin("MyBot v1.0"),
        dashgram.WithHTTPClient(httpClient),
        dashgram.WithUseAsync(),
        dashgram.WithNumWorkers(5),
    )
    defer client.Close()

    // Track events asynchronously
    for i := 0; i < 10; i++ {
        event := map[string]interface{}{
            "user_id": 123456789 + i,
            "action":  "button_clicked",
            "button":  "start",
        }
        client.TrackEventAsync(event)
    }
}
```

## Integration Examples

### With `tgbotapi`

```go
package main

import (
    "log"
    "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/Dashgram/go-dashgram"
)

func main() {
    // Initialize Telegram bot
    bot, err := tgbotapi.NewBotAPI("your-telegram-bot-token")
    if err != nil {
        log.Fatal(err)
    }

    // Initialize Dashgram client
    dashgramClient := dashgram.New(12345, "your-dashgram-access-key")
    defer dashgramClient.Close()

    // Set up update channel
    updateConfig := tgbotapi.NewUpdate(0)
    updateConfig.Timeout = 60
    updates := bot.GetUpdatesChan(updateConfig)

    for update := range updates {
        dashgramClient.TrackEventAsync(update)
        if update.Message != nil {
            if update.Message.IsCommand() {
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Command received!")
                bot.Send(msg)
            }
        }
    }
}
```

### With `telebot`

```go
package main

import (
    "log"
    "time"
    "gopkg.in/telebot.v3"
    "github.com/Dashgram/go-dashgram"
)

func main() {
    bot, err := telebot.NewBot(telebot.Settings{
        Token:  "your-telegram-bot-token",
        Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
    })
    if err != nil {
        log.Fatal(err)
    }

    dashgramClient := dashgram.New(12345, "your-dashgram-access-key")
    defer dashgramClient.Close()

    bot.Handle("/start", func(c telebot.Context) error {
        dashgramClient.TrackEventAsync(c.Update())

        return c.Send("Welcome! I'm your bot.")
    })

    bot.Handle(telebot.OnText, func(c telebot.Context) error {
        dashgramClient.TrackEventAsync(c.Update())

        return c.Send("I received your message!")
    })

    bot.Start()
}
```

## API Reference

### Client Creation

```go
// Basic client
client := dashgram.New(projectID, accessKey)

// With options
client := dashgram.New(projectID, accessKey,
    dashgram.WithAPIURL("https://custom-api.dashgram.io/v1"),
    dashgram.WithOrigin("MyBot v1.0"),
    dashgram.WithHTTPClient(customHTTPClient),
    dashgram.WithUseAsync(),
    dashgram.WithNumWorkers(5),
)
```

### Available Options

- `WithAPIURL(url string)`: Set custom API URL
- `WithOrigin(origin string)`: Set custom origin string
- `WithHTTPClient(client HttpClient)`: Set custom HTTP client
- `WithUseAsync()`: Enable asynchronous processing by default  (client.TrackEvent(...) will act as client.TrackEventAsync(...))
- `WithNumWorkers(num int)`: Set number of worker goroutines to process async events

### Methods

#### Synchronous Methods

```go
// Track an event
err := client.TrackEvent(event)

// Track an event with context
err := client.TrackEventWithContext(ctx, event)

// Track user invitation
err := client.InvitedBy(userID, invitedBy)

// Track user invitation with context
err := client.InvitedByWithContext(ctx, userID, invitedBy)
```

#### Asynchronous Methods

```go
// Track an event asynchronously
client.TrackEventAsync(event)

// Track an event asynchronously with context
client.TrackEventAsyncWithContext(ctx, event)

// Track user invitation asynchronously
client.InvitedByAsync(userID, invitedBy)

// Track user invitation asynchronously with context
client.InvitedByAsyncWithContext(ctx, userID, invitedBy)
```

### Error Handling

```go
import "github.com/Dashgram/go-dashgram"

// Check for specific error types
if err := client.TrackEvent(event); err != nil {
    switch e := err.(type) {
    case *dashgram.InvalidCredentialsError:
        log.Printf("Invalid credentials: %v", e)
    case *dashgram.DashgramAPIError:
        log.Printf("API error (status %d): %s", e.StatusCode, e.Details)
    default:
        log.Printf("Unexpected error: %v", e)
    }
}
```

## Best Practices

1. **Use Async for High-Volume**: Enable async processing for bots with high message volumes
2. **Include Context**: Use context-aware methods for better control over request lifecycle
3. **Handle Errors**: Always check for errors and handle them appropriately
4. **Close Client**: Always call `client.Close()` when shutting down your application
5. **Structured Events**: Use telegram native updates type for better analytics

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
