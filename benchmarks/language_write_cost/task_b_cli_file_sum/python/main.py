try:
    with open("names.txt") as f:
        content = f.read()
except FileNotFoundError:
    print("Error: could not read names.txt")
else:
    for line in content.split("\n"):
        if line != "":
            print("Hello,", line)
        else:
            print()
