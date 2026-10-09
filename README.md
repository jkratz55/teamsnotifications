# teamsnotifications

A very simple zero dependency library for posting JSON payloads to Microsoft Teams through webhooks. It includes a few basic types for Adaptive Cards, but does not aim to support every Adaptive Card feature. This covers basic use cases such as sending notifications in a CI/CD pipeline, from an application, Kubernetes controller, or some other automation component. The caller is responsible for providing the payload format required by their webhook endpoint. The main reason I published this as a library is I often find myself needing this feature and copy-pasting the code between code bases.

## Usage

The library provides two ways of posting payloads to Teams.

1. Creating a `Client` and using `PostPayload` on the `Client`
2. Using the `PostPayload` function in the package

There is very little difference between the two methods, so which one you use will often depend on personal preference. The `Client` is initialized with the webhook, so you do not need to pass the webhook as an argument when calling `PostPayload` on the `Client` type. The `Client` type also takes options if you want to customize the HTTP client. Meanwhile, the package-level `PostPayload` function requires you to pass the webhook each time and always uses `http.DefaultClient` to send the request. Both methods accept any JSON-serializable payload.
