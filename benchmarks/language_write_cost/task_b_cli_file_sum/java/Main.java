import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;

public class Main {
    public static void main(String[] args) {
        String content;
        try {
            content = new String(Files.readAllBytes(Paths.get("names.txt")));
        } catch (IOException e) {
            System.out.println("Error: could not read names.txt");
            return;
        }

        for (String line : content.split("\n")) {
            if (!line.equals("")) {
                System.out.println("Hello, " + line);
            } else {
                System.out.println("");
            }
        }
    }
}
