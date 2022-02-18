BODY=`curl --silent -X POST "http://localhost:4000/openrtb" --header "Content-Type: application/json" --data-raw '{"id":"5c270986-b6b0-48a5-8a1f-6c8ef5eb18f3","at":1,"device":{"ip":"72.66.34.85","geo":{"type":2,"ipservice":3,"city":"X0-Excluded"},"ua":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36","lmt":0,"ifa":"00000000-0000-0000-0000-000000000000","language":"en","os":"Web","ext":{"has_adblock":0},"devicetype":2},"user":{"id":"3d03989298ed630fb94d17927ec35b142dc3d58d6b3842b9e8c17f0f512c50d6","yob":1933,"data":[]},"regs":{"coppa":0,"ext":{"gdpr":0,"us_privacy":"1---"}},"site":{"id":"987654321","domain":"tumblr.com","page":"https://www.tumblr.com"},"imp":[{"id":"1","tagid":"3","secure":1}]}'`

if [ ! -z $BODY ]
then
    IMPRESSION_BEACON=`echo $BODY | jq -r '.seatbid[0].bid[0].adm | fromjson | .native.eventtrackers[0].url'`
    curl --silent $IMPRESSION_BEACON &
    echo "Beacon fired"
else
    echo "Empty response"
fi
