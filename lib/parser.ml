open Token

module Parser = struct
  type context = { current : int; tokens : Token.token list }

  let advance parser = None
  let start tokens = None
end
