string content;
try
{
    content = File.ReadAllText("names.txt");
}
catch (FileNotFoundException)
{
    Console.WriteLine("Error: could not read names.txt");
    return;
}

foreach (var line in content.Split("\n"))
{
    if (line != "")
    {
        Console.WriteLine("Hello, " + line);
    }
    else
    {
        Console.WriteLine("");
    }
}
