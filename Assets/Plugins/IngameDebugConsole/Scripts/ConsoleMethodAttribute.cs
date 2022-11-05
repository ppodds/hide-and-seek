using System;

namespace IngameDebugConsole
{
    [AttributeUsage(AttributeTargets.Method, Inherited = false, AllowMultiple = true)]
    public class ConsoleMethodAttribute : Attribute
    {
        public ConsoleMethodAttribute(string command, string description, params string[] parameterNames)
        {
            Command = command;
            Description = description;
            ParameterNames = parameterNames;
        }

        public string Command { get; }

        public string Description { get; }

        public string[] ParameterNames { get; }
    }
}