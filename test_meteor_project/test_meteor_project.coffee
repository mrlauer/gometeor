if Meteor.is_client
  Session.set 'hello', 'ohai!'
  Template.hello.greeting = ()->
    return "#{Session.get('hello')}."

  m = Meteor.connect "127.0.0.1:3010"
  m.methods
    'Greeting' : (name)->"Stub"
    'Bogus' : (name)->"Stub"

  Template.hello.events =
    'click #good' : ()->
        r = m.call 'Greeting', 'Honored guest', (error, result) ->
            if error
                Session.set 'hello': "Got an error: #{error}"
            else
                Session.set 'hello', result
        Session.set 'hello', r

    'click #bad' : ()->
        r = m.call 'Bogus', 'Loser', (error, result) ->
            if error
                Session.set 'hello', "Got an error: #{error}"
            else
                Session.set 'hello', result
        Session.set 'hello', r
