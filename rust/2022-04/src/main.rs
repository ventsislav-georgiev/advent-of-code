use aoc::exec;

fn main() {
    exec(&[task1, task2]);
}

fn task1(input: String) {
    let mut result = 0;
    for line in input.lines() {
        let pairs: Vec<&str> = line.split(",").collect();
        let a: Vec<&str> = pairs[0].split("-").collect();
        let b: Vec<&str> = pairs[1].split("-").collect();
        let a_x = a[0].parse::<i32>().unwrap();
        let a_y = a[1].parse::<i32>().unwrap();
        let b_x = b[0].parse::<i32>().unwrap();
        let b_y = b[1].parse::<i32>().unwrap();
        if (a_x <= b_x && a_y >= b_y) || (b_x <= a_x && b_y >= a_y) {
            result += 1;
        }
    }

    println!("{}", result);
}

fn task2(input: String) {
    let mut result = 0;
    for line in input.lines() {
        let pairs: Vec<&str> = line.split(",").collect();
        let a: Vec<&str> = pairs[0].split("-").collect();
        let b: Vec<&str> = pairs[1].split("-").collect();
        let a_x = a[0].parse::<i32>().unwrap();
        let a_y = a[1].parse::<i32>().unwrap();
        let b_x = b[0].parse::<i32>().unwrap();
        let b_y = b[1].parse::<i32>().unwrap();
        if (a_x <= b_x && a_y >= b_x) || (b_x <= a_x && b_y >= a_x) {
            result += 1;
        }
    }

    println!("{}", result);
}
