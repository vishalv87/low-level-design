def remove_words_with_digits(input_file, output_file):
    with open(input_file, 'r') as f:
        words = f.readlines()

    # Filter out words containing any digits
    clean_words = [word.strip() for word in words if not any(char.isdigit() for char in word)]

    with open(output_file, 'w') as f:
        for word in clean_words:
            f.write(word + '\n')

    print(f"Cleaned {len(words) - len(clean_words)} words with digits. Output written to {output_file}")


# Example usage
remove_words_with_digits("oxford_3000.txt", "oxford3000_clean.txt")
