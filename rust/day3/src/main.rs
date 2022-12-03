use aoc::exec;

fn main() {
    exec(&[task1, task2]);
}

fn task1(input: String) {
    let mut result = 0;
    for r in input.lines() {
        let halflen = r.len() / 2;
        let bitmap = to_bitmap(&r[..halflen]) & to_bitmap(&r[halflen..]);
        for p in 0..53 {
            if bitmap & (1 << p) != 0 {
                result += p + 1;
            }
        }
    }

    println!("{}", result);
}

fn task2(input: String) {
    let lines = input.lines().collect::<Vec<&str>>();
    let mut result = 0;

    for i in (0..lines.len()).step_by(3) {
        let r1 = lines[i];
        let r2 = lines[i + 1];
        let r3 = lines[i + 2];

        let bitmap = to_bitmap(&r1) & to_bitmap(&r2) & to_bitmap(&r3);
        for p in 0..53 {
            if bitmap & (1 << p) != 0 {
                result += p + 1;
            }
        }
    }

    println!("{}", result);
}

fn to_bitmap(s: &str) -> u64 {
    let low_a = 'a' as u64;
    let high_a = 'A' as u64;
    let mut bitmap = 0;
    for ch in s.chars() {
        if ch >= 'a' {
            bitmap |= 1 << (ch as u64 - low_a);
        } else {
            bitmap |= 1 << (ch as u64 - high_a + 26);
        }
    }
    return bitmap;
}
