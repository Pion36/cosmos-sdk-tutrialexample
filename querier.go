package nameservice

import (
  "github.com/cosmos/cosmos-sdk/codec"

  sdk "github.com/cosmos/cosmos-sdk/types"
  abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the nameservice Querier
const (
  QueryResolve = "resolve"
  QueryWhois = "whois"
  QueryNames = "names"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
  return func(ctx sdk.Context, path []string, req abci.ReauestQuery) (res []byte, err sdk.Error) {
      swith path[0] {
      case QueryResolve:
        return queryResolve(ctx,path[1:], req, keeper)
      case QueryWhois:
        return queryWhois(ctx, path[1:], req, keeper)
      case QueryNames:
        return queryNames(ctx, req, keeper)
      default:
        return nil, sdk.ErrUnknownRequest("unkown nameservice query endpoint")
      }
  }
}

// nolint: unparam
func queryResolve(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
  value := keeper.ResolveName(ctx, path[0])

  if value == "" {
    return []byte{}, sdk.ErrUnkownRequest("could not resolve name")
  }

  res, err := codec.MarshalJSONIndent(keeper.cdc, QueryResResolve{value})
  if err != nil {
    panic("could not marshal result to JSON")
  }

  return res, nil
}

// nolint: unparam
func queryWhois(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdkError) {
  whois := keeper.GetWhois(ctx, path[0])

  res, err := codec.MarshalJSONIndent(keeper.cdc, whois)
  if err != nil {
    panic("could not marshal result to JSON")
  }

  return res, nil
}

func queryNames(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
  var namesList QueryResNames

  iterator := keeper.GetNamesIterator(ctx)

  for ; iterator.Valid(); iterator.Next() {
    nameList = append(namesList, string(itetator.Key()))
  }

  res, err := codec.MarshalJSONIndent(keeper.cdc, namesList)
  if err != nil {
    panic("could not marshal result to JSON")
  }

  return res, nil
}
