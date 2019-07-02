require_relative 'placeholders'

require 'deepsort'
require 'json'

step "metrics file :path has permissions :permissions" do |path, permissions|
  expect(File.file?(path)).to be(true)

  actual = File.stat(path).mode.to_s(8).split('')[-4..-1].join
  expect(actual).to eq(permissions)
end

step "metrics file :path should have following keys:" do |path, keys|
  expected_keys = keys.split("\n").map(&:strip).reject { |x| x.empty? }

  eventually(timeout: 10) {
    expect(File.file?(path)).to be(true), "path #{path} is not file"
  }

  metrics_keys = File.open(path, 'rb') { |f| JSON.parse(f.read).keys }

  expect(metrics_keys).to match_array(expected_keys)
end
