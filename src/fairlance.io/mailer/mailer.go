package mailer

type Mailer interface {
    SendWelcomeMessage(string) (string, error)
}

const WelcomeMessage = `
Hello,

Welcome to the Fairlance community.
We would like to build Fairlance as a community work platform dedicated to establishing a new business paradigm based on principles of responsibility and fairness.

But who gets to decide what is fair and what is not….?

Well... all of us!

Therefore we need your thoughts on how to make this system work best for all, your experiences of both good and bad freelance practices, publicly shared and discussed. As we believe that strong communication between freelancers, clients and platform is crucial.

So let's build it!

We are currently working on BETA version and gathering feedback from the get go. Everyone is invited to bring in new ideas and contribute to our efforts.

We will never spam or give out your email and you can unsubscribe at any point without any hassle. We will only send you emails when there are new and exciting updates coming.

Join us on/Talk to us on(Facebook, Linkedin, Twitter, Reddit, )

We would love to hear from you soon,

Fairlance team
`
