import org.junit.jupiter.api.Test;
import static org.junit.jupiter.api.Assertions.assertEquals;

public class CalculatorTests {
    @Test
    void addReturnsCorrectSum() {
        assertEquals(5, Calculator.add(2, 3));
    }
}
