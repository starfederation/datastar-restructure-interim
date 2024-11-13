module internal StarFederation.Datastar.Utility

open System
open Microsoft.FSharp.Reflection

let unionCaseFromString<'a> (str:string) args =
    match FSharpType.GetUnionCases(typeof<'a>) |> Array.filter (fun unionCaseInfo -> unionCaseInfo.Name.ToLower() = str.ToLower()) with
    | [| unionCaseInfo |] -> ValueSome (FSharpValue.MakeUnion( unionCaseInfo, args ) :?> 'a)
    | _ -> ValueNone

let lowerFirstCharacter item =
    String.Concat(Char.ToLowerInvariant($"%A{item}"[0]), $"%A{item}".Substring(1))
