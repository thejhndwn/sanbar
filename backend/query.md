UPDATE solo_survival_games
SET
    solve_timestamps = array_append(solve_timestamps, NOW()),
    game_index = game_index + 1,
    scores = array_append(scores, $1),  -- score for this solve
    score = score + $1,
    skips = skips + $2,  -- if skip was used
    skips_indices = CASE WHEN $2 > 0 THEN array_append(skips_indices, game_index) ELSE skips_indices END,

    -- Compute time spent on this problem
    time_spent = array_append(
        time_spent,
        ROUND(
            EXTRACT(EPOCH FROM (
                NOW() - 
                -- Use last solve time, or start_time if first problem
                COALESCE(
                    solve_timestamps[ARRAY_LENGTH(solve_timestamps, 1)],  -- last solve
                    start_time                                                  -- fallback
                )
            ))
        )
    ),

    -- Reduce remaining `time` by time spent
    time = GREATEST(
        0,
        time - ROUND(
            EXTRACT(EPOCH FROM (
                NOW() - 
                COALESCE(
                    solve_timestamps[ARRAY_LENGTH(solve_timestamps, 1)],
                    start_time
                )
            ))
        )
    ),

    -- Bonus: Add time on solve (e.g., +10s)
    time = time + 10,

    -- Subtract time on skip (e.g., -20s)
    time = time - ($2 * 20),

    updated_at = NOW()
WHERE id = $3
RETURNING 
    combos, 
    game_index, 
    score, 
    time,
    time_spent,
    ARRAY_LENGTH(combos, 1) - (game_index + 1) AS problems_remaining;
