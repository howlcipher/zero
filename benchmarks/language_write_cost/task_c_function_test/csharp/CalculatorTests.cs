using Xunit;

public class CalculatorTests
{
    [Fact]
    public void AddReturnsCorrectSum()
    {
        Assert.Equal(5, Calculator.Add(2, 3));
    }
}
