use aoc::exec;
use phf::phf_map;

fn main() {
    exec(&[task1, task2]);
}

static POINTS: phf::Map<&str, i32> = phf_map! {
    "A" => 0,
    "B" => 1,
    "C" => 2,
    "X" => 0,
    "Y" => 1,
    "Z" => 2,
};

fn task1(input: String) {
    let mut result = 0;
    for line in input.lines() {
        let f: Vec<&str> = line.split(" ").collect();
        let p1 = POINTS[f[0]];
        let p2 = POINTS[f[1]];

        let winner = (3 + p1 - p2) % 3;
        match winner {
            0 => result += 3 + p2 + 1,
            1 => result += 0 + p2 + 1,
            2 => result += 6 + p2 + 1,
            _ => panic!("Unknown winner"),
        }
    }

    println!("{}", result);
}

fn task2(input: String) {
    let mut result = 0;
    for line in input.lines() {
        let f: Vec<&str> = line.split(" ").collect();
        let p1 = POINTS[f[0]];
        let mut p2 = POINTS[f[1]];

        match p2 {
            0 => p2 = p1 - 1,
            1 => p2 = p1,
            2 => p2 = p1 + 1,
            _ => panic!("Unknown winner"),
        }

        if p2 < 0 {
            p2 = 2
        } else if p2 > 2 {
            p2 = 0
        }

        let winner = (3 + p1 - p2) % 3;
        match winner {
            0 => result += 3 + p2 + 1,
            1 => result += 0 + p2 + 1,
            2 => result += 6 + p2 + 1,
            _ => panic!("Unknown winner"),
        }
    }

    println!("{}", result);
}
