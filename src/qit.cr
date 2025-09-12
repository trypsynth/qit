require "./command_registry"

module Qit
  extend self

  def run
    CommandRegistry.initialize_commands
    command_name = ARGV.empty? ? "help" : ARGV[0].downcase
    args = ARGV.empty? ? [] of String : ARGV[1..]
    if command = CommandRegistry.get(command_name)
      command.execute(args)
      exit 0
    else
      STDERR.puts "Unknown command: #{command_name}."
      STDERR.puts "Use 'qit help' to see available commands."
      exit 1
    end
  end
end

Qit.run
