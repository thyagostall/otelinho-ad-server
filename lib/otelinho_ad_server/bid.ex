defmodule OtelinhoAdServer.Bid do
  alias OtelinhoAdServer.BeaconGenerator

  def bid_response(campaign) do
    payload = BeaconGenerator.generate_payload(campaign)

    %{
      seatbid: [
        %{
          bid: [
            %{
              demand_source: "direct",
              price: Decimal.to_float(campaign.max_bid),
              cid: to_string(campaign.id),
              id: UUID.uuid4(),
              adm: Poison.encode!(adm(campaign, payload)),
              nurl: BeaconGenerator.generate_url("win", payload),
              lurl: BeaconGenerator.generate_url("loss", payload),
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

  defp adm(campaign, beacon_payload) do
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
            url: BeaconGenerator.generate_url("impression", beacon_payload),
            event: 1
          },
          %{
            method: 1,
            url: BeaconGenerator.generate_url("%{EVENT_TYPE}", beacon_payload),
            event: 600
          }
        ],
        ver: "1.2"
      }
    }
  end
end
