ALTER TABLE auth.usu_usuario
ADD COLUMN usu_tx_refresh_token varchar(100),
ADD COLUMN usu_dt_refresh_token_exp TIMESTAMPTZ
