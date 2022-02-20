defmodule OtelinhoAdServerWeb.OpenRTB do
  use OtelinhoAdServerWeb, :controller

  alias OtelinhoAdServer.Index
  alias OtelinhoAdServer.Auction
  alias OtelinhoAdServer.Pacing
  alias OtelinhoAdServer.Bid

  def openrtb(conn, _params) do
    Index.retrieve_active_campaigns()
    |> Pacing.filter()
    |> Auction.run_auction()
    |> bid_response(conn)
  end

  defp bid_response(nil, conn), do: send_resp(conn, 204, "")
  defp bid_response(campaign, conn), do: json(conn, Bid.bid_response(campaign))
end
