# TCP-Messaging
Implements a basic tcp messaging system with one server and multiple clients


### Time Taken
 - Preliminary Reading (2 hrs)
 - Implementation (11 hrs)

### Note
 - All commands are delimited by '~'
 - Development branch was not used since it was not required here
 - golangci-lint has been used and run though no .yml file has been attached

### Sample Operations
  - {"operation":"ClientList","value":""}~

  - {"operation":"ForwardMessage","value":"-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEoqHXQ+0LGt1idkyME7AJ6cpKrjaY\ns4rJdWZTh9dpzKkfssuC19g610SeampqDZMY6HIbNJNysitIwsiyZ6/gyw==\n-----END PUBLIC KEY-----\n","target":"-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAELY1A6HRi34vACLj2TqzBUdQY73hA\nFrJwxI5Bydc4wxK9YG91jXQEHA++fdGkvdRIs07CShpuSb98iJxy5S4Yog==\n-----END PUBLIC KEY-----\n"}~
