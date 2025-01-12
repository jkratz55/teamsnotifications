# teamsnotifications

A very simple zero dependency library for posting basic messages to Microsoft Teams through webhooks. This library does **not** support AdaptiveCards and probably never will. This covers basic use cases such as sending notifications in a CI/CD pipeline, from an application, Kubernetes controller, or some other automation component. Its goals are not to support building fancy features supported by Adaptive Cards. The main reason I published this as a library is I often find myself needing this feature and copy-pasting the code between code bases.

## Usage

The library provides two ways of posting messages/notifications to teams.

1. Creating a `Client` and using `PostMessage` on the `Client`
2. Using `PostMessage` function in the package

There is very little difference between the two methods so which one you use will often depend on personal preference. The `Client` is initialized with the webhook, thus you won't need to pass the webhook as an argument when calling `PostMessage` on the `Client` type. The `Client` type also takes options if you want to customize the HttpClient. Meanwhile, the `PostMessage` function in the package requires you pass the webhook each time, and it always uses `http.DefaultClient` to send the request.

