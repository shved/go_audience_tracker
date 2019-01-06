`curl --get localhost:9292/pulse -d customer_id=1 -d video_id=1`
sleep(0.3)
`curl --get localhost:9292/pulse -d customer_id=1 -d video_id=2`
sleep(0.3)
`curl --get localhost:9292/pulse -d customer_id=1 -d video_id=3`
sleep(0.3)
`curl --get localhost:9292/pulse -d customer_id=2 -d video_id=1`
sleep(0.3)
`curl --get localhost:9292/customers/1`
sleep(0.3)
`curl --get localhost:9292/videos/1`
sleep(0.3)
`curl --get localhost:9292/videos/1`
sleep(0.3)
`curl --get localhost:9292/videos/1`
sleep(0.3)
`curl --get localhost:9292/videos/1`
sleep(0.3)
