defmodule OtelinhoAdServerWeb.OpenRTB do
  use OtelinhoAdServerWeb, :controller

  def openrtb(conn, _params) do
    json(conn, bid_response())
  end

  defp bid_response() do
    %{
      seatbid: [
        %{
          bid: [
            %{
              demand_source: "direct",
              price: 7,
              cid: "1",
              id: "3c8e88f7-9be3-46c3-8c83-26a69fd68e6d",
              adm: "{\"native\":{\"assets\":[{\"id\":1,\"data\":{\"type\":501,\"value\":\"t:pWE-0YwL2ycRagbqsCSBuQ;642229909710946305\"},\"required\":1}],\"eventtrackers\":[{\"method\":1,\"url\":\"http://localhost:4000/event/impression/data\",\"event\":1},{\"method\":1,\"url\":\"http://localhost:4000/event/${EVENT_TYPE}/data\",\"event\":600}],\"ver\":\"1.2\"}}",
              nurl: "http://localhost:4000/winnotice",
              adomain: [
                  ""
              ],
              cat: [
                  "IAB12-3"
              ],
              crid: "1",
              impid: "25eed2e8-6520-47cb-a22c-15ef9b6af4c1",
              adid: "1",
              adm_media_type: "native"
            }
          ],
          seat: "1"
        }
      ],
      id: "1"
    }
  end
end
