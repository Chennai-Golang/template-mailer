#!/usr/bin/env ruby

require 'uri'

def add_params(url, params = {})
  uri = URI(url)
  params    = Hash[URI.decode_www_form(uri.query || '')].merge(params)
  uri.query =      URI.encode_www_form(params)
  uri.to_s
end

tweet_link = add_params("https://twitter.com/intent/tweet", {
  hashtags: "ChnGoTurns10, GoTurns10",
  url: "https://www.meetup.com/Chennai-golang-Meetup/events/266846525/",
  text: "I am attending this month's Go meetup happening @ Qube Cinemas"
})

puts tweet_link
