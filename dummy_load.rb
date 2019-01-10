threads = []
customers = (1..1000).to_a
videos = (1..100).to_a

customers.each do |customer|
  threads << Thread.new do
    duration = (5..15).to_a.sample
    video = videos.sample

    # starting offset
    sleep([1, 2, 3, 4, 6, 7, 8, 9, 11, 12, 13, 14].sample)

    duration.times do
      `curl --get -s localhost:9292/pulse -d customer_id=#{customer} -d video_id=#{video}`
      sleep 5
    end
  end
end

threads << Thread.new do
  20.times do
    interval = (3..7).to_a.sample
    video = videos.sample
    sleep interval
    `curl --get -s localhost:9292/videos/#{video}`
  end
end

threads << Thread.new do
  20.times do
    interval = (3..7).to_a.sample
    customer = customers.sample
    sleep interval
    `curl --get -s localhost:9292/customers/#{customer}`
  end
end

threads.each(&:join)
