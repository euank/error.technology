require_relative './errortech.rb'
require 'json'

10.times do
  begin
  JSON.parse("asdf")
  rescue Exception => e
    puts e
  end
end
