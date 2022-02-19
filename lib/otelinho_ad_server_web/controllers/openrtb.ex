defmodule OtelinhoAdServerWeb.OpenRTB do
  use OtelinhoAdServerWeb, :controller

  alias OtelinhoAdServer.Index
  alias OtelinhoAdServer.Auction

  def openrtb(conn, _params) do
    response = Index.retrieve_active_campaigns()
    |> Auction.run_auction()
    |> bid_response()

    json(conn, response)
  end

  defp bid_response(campaign) do
    %{
      seatbid: [
        %{
          bid: [
            %{
              demand_source: "direct",
              price: Decimal.to_float(campaign.max_bid),
              cid: to_string(campaign.id),
              id: UUID.uuid4(),
              adm: Poison.encode!(adm(campaign)),
              nurl: "http://localhost:4000/win",
              lurl: "http://localhost:4000/loss",
              adomain: [
                  ""
              ],
              cat: [
                  "IAB12-3"
              ],
              crid: "1",
              impid: UUID.uuid4(),
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

  defp adm(campaign) do
    %{
      native: %{
        assets: [
          %{
            id: 1,
            data: %{
              type: 501,
              value: campaign.creative
            },
            required: 1
          }
        ],
        eventtrackers: [
          %{
            method: 1,
            url: "http://localhost:4000/event/impression/data",
            event: 1
          },
          %{
            method: 1,
            url: "http://localhost:4000/event/${EVENT_TYPE}/data",
            event: 600
          }
        ],
        ver: "1.2"
      }
    }
  end
end
