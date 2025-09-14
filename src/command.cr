abstract class Command
  abstract def name : String
  abstract def description : String
  abstract def execute(args : Array(String)) : Nil
end
