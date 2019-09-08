use std::env;

pub fn fizzbuzz(n: i32) -> String {
    match (n % 3, n % 5) {
        (0, 0) => "FizzBuzz".to_string(),
        (0, _) => "Fizz".to_string(),
        (_, 0) => "Buzz".to_string(),
        (_, _) => n.to_string(),
    }
}


#[cfg(test)]
mod test {
    use super::*;
    #[test]
    fn test_fizzbuzz() {
        assert_eq!(fizzbuzz(0), "FizzBuzz");
        assert_eq!(fizzbuzz(1), "1");
        assert_eq!(fizzbuzz(2), "2");
        assert_eq!(fizzbuzz(3), "Fizz");
        assert_eq!(fizzbuzz(4), "4");
        assert_eq!(fizzbuzz(5), "Buzz");
        assert_eq!(fizzbuzz(6), "Fizz");

        assert_eq!(fizzbuzz(14), "14");
        assert_eq!(fizzbuzz(15), "FizzBuzz");
        assert_eq!(fizzbuzz(16), "16");
    }
}

fn run(args: Vec<String>) -> i32 {
    let arg = match args.len() {
        0 | 1 => None,
        _ => args[1].parse::<i32>().ok()
    };
    match arg {
        Some(arg) => {
            for x in 1..(arg+1) {
                println!("{}", fizzbuzz(x));
            }
            0
        },
        None => {
            println!("Usage: {} number", args[0]); 1
        },
    }
}

fn main() {
    let args : Vec<String> = env::args().collect();
    ::std::process::exit(run(args));
}
