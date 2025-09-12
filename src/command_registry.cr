require "./command"
require "./commands/*"

module Qit
  class CommandRegistry
    @@commands = {} of String => Command

    def self.register(command : Command) : Nil
      @@commands[command.name] = command
    end

    def self.get(name : String) : Command?
      @@commands[name]?
    end

    def self.all : Hash(String, Command)
      @@commands
    end

    def self.initialize_commands : Nil
      [
        Commands::AcpCommand.new,
        Commands::AmendCommand.new,
        Commands::CpCommand.new,
        Commands::DbCommand.new,
        Commands::HelpCommand.new,
        Commands::IgnoreCommand.new,
        Commands::LastCommand.new,
        Commands::LicenseCommand.new,
        Commands::LogCommand.new,
        Commands::NbCommand.new,
        Commands::NewCommand.new,
        Commands::ResetCommand.new,
        Commands::StatusCommand.new,
        Commands::UndoCommand.new,
      ].each { |command| register(command) }
    end
  end
end
