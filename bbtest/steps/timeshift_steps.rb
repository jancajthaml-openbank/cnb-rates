require "date"

step "current time is :timeshift" do |timeshift|
  @timeshift = DateTime.parse(timeshift)
end
