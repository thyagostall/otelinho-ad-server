defmodule OtelinhoAdServerWeb.OpenRTB do
  use OtelinhoAdServerWeb, :controller

  alias OtelinhoAdServer.Index
  alias OtelinhoAdServer.Auction
  alias OtelinhoAdServer.BeaconGenerator
  alias OtelinhoAdServer.Pacing

  def openrtb(conn, _params) do
    campaign = Index.retrieve_active_campaigns()
    |> Pacing.filter()
    |> Auction.run_auction()

    case campaign do
      nil -> send_resp(conn, 204, "")
      _ -> json(conn, bid_response(campaign))
    end
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
              nurl: BeaconGenerator.generate("win", campaign),
              lurl: BeaconGenerator.generate("loss", campaign),
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
            url: BeaconGenerator.generate("impression", campaign),
            event: 1
          },
          %{
            method: 1,
            url: BeaconGenerator.generate("%{EVENT_TYPE}", campaign),
            event: 600
          }
        ],
        ver: "1.2"
      }
    }
  end
end
