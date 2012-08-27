if Meteor.is_client
  Session.set 'hello', 'ohai!'
  Template.hello.greeting = ()->
    return "#{Session.get('hello')}."

  m = Meteor.connect "127.0.0.1:3010"
  m.methods
    'Greeting' : (name)->"Stub"
    'Bogus' : (name)->"Stub"

  fnCallback = (error, result) ->
        if error
            Session.set 'hello', "Got an error: #{error.error} #{error.reason}"
        else
            Session.set 'hello', result

  Template.hello.events =
    'click #good' : ()->
        r = m.call 'Greeting', 'Honored guest', fnCallback
        Session.set 'hello', r

    'click #complicated' : ()->
        r = m.call 'Complicated', 'Honored guest', { Age: 137, Description: 'Blah' }, fnCallback
        Session.set 'hello', r

    'click #bad_args' : ()->
        r = m.call 'Complicated', 'Honored guest', { foo: 42, bar: 'Blah' }, 'foo', fnCallback
        Session.set 'hello', r

    'click #bad' : ()->
        r = m.call 'Bogus', 'Loser', fnCallback
        Session.set 'hello', r
