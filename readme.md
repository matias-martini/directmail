# DirectMail: Direct Mail Delivery in Go
## Introduction

Welcome to DirectMail, the Go library that takes the "roundabout" out of email delivery. We're not just cutting corners; we're slicing through the entire SMTP server maze. Get ready to send emails in Go as if you're just shouting across the room - direct, a bit loud, and no looking back.

## Why DirectMail?
In the labyrinth of complex email delivery systems, DirectMail is that rare straight line.
It's like a unicorn, but less sparkly and more... email-y. We realized the world didn't need another Swiss Army knife of email delivery; what it needed was a good old-fashioned hammer. And here we are.

## Example Usage
Here's a basic example of how to send an email using DirectMail:

```go
package main

import (
    "github.com/matias-martini/directmail"
)

func main() {
    directmail.SendEmail(
        "sender@yourdomain.com",
        "recipient@theirdomain.com",
        "Here goes the subject",
        "Here goes the body"
    )
}
```

## Key Features
- Direct to the Point: We cut out the middleman (SMTP servers) and connect you straight to the recipient's server.
- Keeping it Light: We're all about keeping things simple. No frills, no unnecessary extras.
- DirectMail is ideal for applications where direct and minimalistic email delivery is required, such as automated system notifications, basic transactional emails, and small-scale mailing needs.

## Installation
To use DirectMail, simply import it into your Go project:

```go
import "github.com/matias-martini/directmail"
```

## Configuration
Setting up DirectMail involves some DNS setup and careful management of a crucial environment variable. Think of it as laying down the tracks for your email train to run smoothly and securely.

### DNS Setup: Your Roadmap to Success

#### DomainKeys Identified Mail (DKIM) Setup

DKIM adds a digital signature to your emails, ensuring they’re not tampered with and verifying they come from your domain.

  - Generate Your Keys: Use OpenSSL to create your private and public key pair.
    ```bash
    openssl genrsa -out private.key 1024
    openssl rsa -in private.key -out public.key -pubout -outform PEM
    ```
    Treat your private.key like a top-secret document.

  - Publish Your Public Key: Add a TXT record in your DNS:
    ```plaintext
    Host: default._domainkey.yourdomain.com
    Value: v=DKIM1; k=rsa; p=[Your public key here without line breaks]
    ```
    This is like hanging your public key on the digital notice board for all to verify.

#### Sender Policy Framework (SPF) Setup

SPF prevents spammers from sending messages with forged from addresses at your domain.

  - Create an SPF Record: This TXT record in your DNS settings tells the world which mail servers are allowed to send emails on behalf of your domain.
    ```
    v=spf1 include:yourdomain.com ~all
    ```
    The `~all` part is like saying, "These are my official mail carriers. If it's not them, be suspicious."

#### Domain-based Message Authentication, Reporting & Conformance (DMARC) Setup

DMARC uses SPF and DKIM to improve email security and helps email receivers determine what to do with emails that don’t pass authentication.

  - Add a DMARC Record: Another TXT record for your DNS settings, it's like setting rules for what to do with emails that pretend to be from your domain.
    ```
    v=DMARC1; p=reject; rua=mailto:reportmail@mail.com
    ```
    The `p=none` policy tells email providers to not enforce any action on unauthenticated emails (for starters). As you get more confident, you can change this to a stricter policy.

### Environment Variable: The Secret Ingredient

Store your private.key generated earlier in an environment variable named DIRECTEMAIL_PRIVATE_KEY.
The private key is the cornerstone of the DKIM signing process. DirectMail uses this key to digitally sign all outgoing emails. The recipient's mail server will then use the corresponding public key, which you have published in the DKIM DNS record, to verify the authenticity of each email.

The environment variable keeps your private key secure and easily accessible for DirectMail without hardcoding it.

Depending on your operating system or deployment environment, set the variable like so:

```
export DIRECTEMAIL_PRIVATE_KEY=$(cat path/to/private.key)
```

And there you have it, the trifecta of DNS records and the secret sauce of environment variables. It's like setting up the ultimate security system for your emails. A bit of effort now for a lot of peace of mind later!

## The Good, the Bad, and the DirectMail
- Pros: You don't need any email server or service to send your emails. You get complete control, and it's lighter than a feather.
- Downsides: DNS setup is mandatory - no shortcuts here, folks. Also we don't retry if the delivery fails. It's like asking someone out, getting a 'no', and just walking away.

Super Basic: If you want bells and whistles, go buy a bike. This is bare-bones by design.


## License
DirectMail is released under the MIT License. See the LICENSE.txt file for more details.

