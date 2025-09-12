require "../command"

module Qit::Commands
  class HelpCommand < Command
    def name : String
      "help"
    end

    def description : String
      "show this help message"
    end

    def execute(args : Array(String)) : Nil
      puts <<-USAGE
      Qit - Quin's tiny Git helper.
      Usage: qit <command> [<args>...]
      Available commands:
      USAGE
      Qit::CommandRegistry.all.each do |name, command|
        puts "  #{name}: #{command.description}."
      end
    end
  end
end
